// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: proto/oidc/idp/idp.proto

/*
Package idp is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package idp

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var (
	_ codes.Code
	_ io.Reader
	_ status.Status
	_ = errors.New
	_ = runtime.String
	_ = utilities.NewDoubleArray
	_ = metadata.Join
)

func request_IDPService_CreateIDP_0(ctx context.Context, marshaler runtime.Marshaler, client IDPServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq CreateIDPRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := client.CreateIDP(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_IDPService_CreateIDP_0(ctx context.Context, marshaler runtime.Marshaler, server IDPServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq CreateIDPRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := server.CreateIDP(ctx, &protoReq)
	return msg, metadata, err
}

func request_IDPService_GetIDP_0(ctx context.Context, marshaler runtime.Marshaler, client IDPServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq GetIDPRequest
		metadata runtime.ServerMetadata
		err      error
	)
	val, ok := pathParams["id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}
	protoReq.Id, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}
	msg, err := client.GetIDP(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_IDPService_GetIDP_0(ctx context.Context, marshaler runtime.Marshaler, server IDPServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq GetIDPRequest
		metadata runtime.ServerMetadata
		err      error
	)
	val, ok := pathParams["id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}
	protoReq.Id, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}
	msg, err := server.GetIDP(ctx, &protoReq)
	return msg, metadata, err
}

func request_IDPService_GetIDPBySlug_0(ctx context.Context, marshaler runtime.Marshaler, client IDPServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq GetIDPBySlugRequest
		metadata runtime.ServerMetadata
		err      error
	)
	val, ok := pathParams["slug"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "slug")
	}
	protoReq.Slug, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "slug", err)
	}
	msg, err := client.GetIDPBySlug(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_IDPService_GetIDPBySlug_0(ctx context.Context, marshaler runtime.Marshaler, server IDPServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq GetIDPBySlugRequest
		metadata runtime.ServerMetadata
		err      error
	)
	val, ok := pathParams["slug"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "slug")
	}
	protoReq.Slug, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "slug", err)
	}
	msg, err := server.GetIDPBySlug(ctx, &protoReq)
	return msg, metadata, err
}

func request_IDPService_DeleteIDP_0(ctx context.Context, marshaler runtime.Marshaler, client IDPServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq DeleteIDPRequest
		metadata runtime.ServerMetadata
		err      error
	)
	val, ok := pathParams["id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}
	protoReq.Id, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}
	msg, err := client.DeleteIDP(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_IDPService_DeleteIDP_0(ctx context.Context, marshaler runtime.Marshaler, server IDPServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq DeleteIDPRequest
		metadata runtime.ServerMetadata
		err      error
	)
	val, ok := pathParams["id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}
	protoReq.Id, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}
	msg, err := server.DeleteIDP(ctx, &protoReq)
	return msg, metadata, err
}

func request_IDPService_UpdateIDP_0(ctx context.Context, marshaler runtime.Marshaler, client IDPServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq UpdateIDPRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := client.UpdateIDP(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_IDPService_UpdateIDP_0(ctx context.Context, marshaler runtime.Marshaler, server IDPServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq UpdateIDPRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := server.UpdateIDP(ctx, &protoReq)
	return msg, metadata, err
}

func request_IDPService_ListIDP_0(ctx context.Context, marshaler runtime.Marshaler, client IDPServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq ListIDPRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := client.ListIDP(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_IDPService_ListIDP_0(ctx context.Context, marshaler runtime.Marshaler, server IDPServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq ListIDPRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := server.ListIDP(ctx, &protoReq)
	return msg, metadata, err
}

// RegisterIDPServiceHandlerServer registers the http handlers for service IDPService to "mux".
// UnaryRPC     :call IDPServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterIDPServiceHandlerFromEndpoint instead.
// GRPC interceptors will not work for this type of registration. To use interceptors, you must use the "runtime.WithMiddlewares" option in the "runtime.NewServeMux" call.
func RegisterIDPServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server IDPServiceServer) error {
	mux.Handle(http.MethodPost, pattern_IDPService_CreateIDP_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/proto.oidc.idp.IDPService/CreateIDP", runtime.WithHTTPPathPattern("/api/v1/idp/_create"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_IDPService_CreateIDP_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_IDPService_CreateIDP_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodGet, pattern_IDPService_GetIDP_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/proto.oidc.idp.IDPService/GetIDP", runtime.WithHTTPPathPattern("/api/v1/idp/id/{id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_IDPService_GetIDP_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_IDPService_GetIDP_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodGet, pattern_IDPService_GetIDPBySlug_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/proto.oidc.idp.IDPService/GetIDPBySlug", runtime.WithHTTPPathPattern("/api/v1/idp/slug/{slug}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_IDPService_GetIDPBySlug_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_IDPService_GetIDPBySlug_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodDelete, pattern_IDPService_DeleteIDP_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/proto.oidc.idp.IDPService/DeleteIDP", runtime.WithHTTPPathPattern("/api/v1/idp/id/{id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_IDPService_DeleteIDP_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_IDPService_DeleteIDP_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodPost, pattern_IDPService_UpdateIDP_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/proto.oidc.idp.IDPService/UpdateIDP", runtime.WithHTTPPathPattern("/api/v1/idp/_update"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_IDPService_UpdateIDP_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_IDPService_UpdateIDP_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodPost, pattern_IDPService_ListIDP_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/proto.oidc.idp.IDPService/ListIDP", runtime.WithHTTPPathPattern("/api/v1/idp/_search"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_IDPService_ListIDP_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_IDPService_ListIDP_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})

	return nil
}

// RegisterIDPServiceHandlerFromEndpoint is same as RegisterIDPServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterIDPServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.NewClient(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()
	return RegisterIDPServiceHandler(ctx, mux, conn)
}

// RegisterIDPServiceHandler registers the http handlers for service IDPService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterIDPServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterIDPServiceHandlerClient(ctx, mux, NewIDPServiceClient(conn))
}

// RegisterIDPServiceHandlerClient registers the http handlers for service IDPService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "IDPServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "IDPServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "IDPServiceClient" to call the correct interceptors. This client ignores the HTTP middlewares.
func RegisterIDPServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client IDPServiceClient) error {
	mux.Handle(http.MethodPost, pattern_IDPService_CreateIDP_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/proto.oidc.idp.IDPService/CreateIDP", runtime.WithHTTPPathPattern("/api/v1/idp/_create"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_IDPService_CreateIDP_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_IDPService_CreateIDP_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodGet, pattern_IDPService_GetIDP_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/proto.oidc.idp.IDPService/GetIDP", runtime.WithHTTPPathPattern("/api/v1/idp/id/{id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_IDPService_GetIDP_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_IDPService_GetIDP_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodGet, pattern_IDPService_GetIDPBySlug_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/proto.oidc.idp.IDPService/GetIDPBySlug", runtime.WithHTTPPathPattern("/api/v1/idp/slug/{slug}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_IDPService_GetIDPBySlug_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_IDPService_GetIDPBySlug_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodDelete, pattern_IDPService_DeleteIDP_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/proto.oidc.idp.IDPService/DeleteIDP", runtime.WithHTTPPathPattern("/api/v1/idp/id/{id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_IDPService_DeleteIDP_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_IDPService_DeleteIDP_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodPost, pattern_IDPService_UpdateIDP_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/proto.oidc.idp.IDPService/UpdateIDP", runtime.WithHTTPPathPattern("/api/v1/idp/_update"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_IDPService_UpdateIDP_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_IDPService_UpdateIDP_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodPost, pattern_IDPService_ListIDP_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/proto.oidc.idp.IDPService/ListIDP", runtime.WithHTTPPathPattern("/api/v1/idp/_search"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_IDPService_ListIDP_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_IDPService_ListIDP_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	return nil
}

var (
	pattern_IDPService_CreateIDP_0    = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "v1", "idp", "_create"}, ""))
	pattern_IDPService_GetIDP_0       = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3, 1, 0, 4, 1, 5, 3}, []string{"api", "v1", "idp", "id"}, ""))
	pattern_IDPService_GetIDPBySlug_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3, 1, 0, 4, 1, 5, 3}, []string{"api", "v1", "idp", "slug"}, ""))
	pattern_IDPService_DeleteIDP_0    = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3, 1, 0, 4, 1, 5, 3}, []string{"api", "v1", "idp", "id"}, ""))
	pattern_IDPService_UpdateIDP_0    = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "v1", "idp", "_update"}, ""))
	pattern_IDPService_ListIDP_0      = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "v1", "idp", "_search"}, ""))
)

var (
	forward_IDPService_CreateIDP_0    = runtime.ForwardResponseMessage
	forward_IDPService_GetIDP_0       = runtime.ForwardResponseMessage
	forward_IDPService_GetIDPBySlug_0 = runtime.ForwardResponseMessage
	forward_IDPService_DeleteIDP_0    = runtime.ForwardResponseMessage
	forward_IDPService_UpdateIDP_0    = runtime.ForwardResponseMessage
	forward_IDPService_ListIDP_0      = runtime.ForwardResponseMessage
)
