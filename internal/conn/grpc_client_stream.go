package conn

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/ydb-platform/ydb-go-genproto/protos/Ydb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ydb-platform/ydb-go-sdk/v3/internal/wrap"
	"github.com/ydb-platform/ydb-go-sdk/v3/internal/xcontext"
	"github.com/ydb-platform/ydb-go-sdk/v3/internal/xerrors"
	"github.com/ydb-platform/ydb-go-sdk/v3/trace"
)

type grpcClientStream struct {
	grpc.ClientStream
	c        *conn
	wrapping bool
	sentMark *modificationMark
	onDone   func(ctx context.Context, md metadata.MD)
	recv     func(error) func(error, trace.ConnState, map[string][]string)
}

func (s *grpcClientStream) CloseSend() (err error) {
	err = s.ClientStream.CloseSend()

	if err != nil {
		if s.wrapping {
			return xerrors.WithStackTrace(
				xerrors.Transport(
					err,
					xerrors.WithAddress(s.c.Address()),
				),
			)
		}
		return xerrors.WithStackTrace(err)
	}

	return nil
}

func (s *grpcClientStream) SendMsg(m interface{}) (err error) {
	cancel := createPinger(s.c)
	defer cancel(xerrors.WithStackTrace(errors.New("send msg finished")))

	err = s.ClientStream.SendMsg(m)

	if err != nil {
		defer func() {
			s.c.onTransportError(s.Context(), err)
		}()

		if s.wrapping {
			err = xerrors.Transport(err,
				xerrors.WithAddress(s.c.Address()),
			)
			if s.sentMark.canRetry() {
				return xerrors.WithStackTrace(xerrors.Retryable(err,
					xerrors.WithName("SendMsg"),
				))
			}
			return xerrors.WithStackTrace(err)
		}

		return err
	}

	return nil
}

func (s *grpcClientStream) RecvMsg(m interface{}) (err error) {
	cancel := createPinger(s.c)
	defer cancel(xerrors.WithStackTrace(errors.New("receive msg finished")))

	defer func() {
		onDone := s.recv(xerrors.HideEOF(err))
		if err != nil {
			md := s.ClientStream.Trailer()
			onDone(xerrors.HideEOF(err), s.c.GetState(), md)
			s.onDone(s.ClientStream.Context(), md)
		}
	}()

	err = s.ClientStream.RecvMsg(m)

	if err != nil {
		defer func() {
			if !xerrors.Is(err, io.EOF) {
				s.c.onTransportError(s.Context(), err)
			}
		}()

		if s.wrapping {
			err = xerrors.Transport(err,
				xerrors.WithAddress(s.c.Address()),
			)
			if s.sentMark.canRetry() {
				return xerrors.WithStackTrace(xerrors.Retryable(err,
					xerrors.WithName("RecvMsg"),
				))
			}
			return xerrors.WithStackTrace(err)
		}

		return err
	}

	if s.wrapping {
		if operation, ok := m.(wrap.StreamOperationResponse); ok {
			if status := operation.GetStatus(); status != Ydb.StatusIds_SUCCESS {
				return xerrors.WithStackTrace(
					xerrors.Operation(
						xerrors.FromOperation(operation),
						xerrors.WithNodeAddress(s.c.Address()),
					),
				)
			}
		}
	}

	return nil
}

func createPinger(c *conn) xcontext.CancelErrFunc {
	c.touchLastUsage()
	ctx, cancel := xcontext.WithErrCancel(context.Background())
	go func() {
		ticker := time.NewTicker(time.Second)
		ctxDone := ctx.Done()
		for {
			select {
			case <-ctxDone:
				ticker.Stop()
				return
			case <-ticker.C:
				c.touchLastUsage()
			}
		}
	}()

	return cancel
}
