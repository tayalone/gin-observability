package nested

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
)

/*Parent - start trancing*/
func Parent(ctx context.Context) {
	tr := otel.Tracer("nested")
	_, span := tr.Start(ctx, "Parent")
	defer span.End()

	sig := make(chan struct{})
	defer close(sig)
	go Child(ctx, 300, sig)
	go Child(ctx, 600, sig)

	<-sig
	<-sig
}

/*Child - internale trancing*/
func Child(ctx context.Context, wait int, sig chan struct{}) {
	tr := otel.Tracer("nested")
	_, span := tr.Start(ctx, "Child")
	go func() {
		time.Sleep(time.Millisecond * time.Duration(wait))
		sig <- struct{}{}
		span.End()
	}()
}
