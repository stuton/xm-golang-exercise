FROM golang:1.20.3-alpine as build

WORKDIR /app

COPY go.mod go.sum ./

# Fetch dependencies.
RUN go mod download

COPY . .

# Build the binaries
RUN CGO_ENABLED=0 go build -o /xm-golang-exercise cmd/api/main.go

# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static-debian11

# Copy our static executable
COPY --from=build /xm-golang-exercise /xm-golang-exercise

# Use an unprivileged user.
USER nonroot:nonroot

CMD ["/xm-golang-exercise"]
