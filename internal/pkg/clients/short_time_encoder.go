package clients

import (
	"time"

	"go.uber.org/zap/zapcore"
)

func ShortTimeEncoder(start time.Time) zapcore.TimeEncoder {
	return func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		offset := t.Sub(start)
		switch {
		case offset > time.Second:
			offset = offset.Round(time.Second / 10)

		case offset > time.Millisecond:
			offset = offset.Round(time.Millisecond / 10)

		case offset > time.Microsecond:
			offset = offset.Round(time.Microsecond / 10)
		}
		pae.AppendString(offset.String())
	}
}
