FROM golang:1.22-alpine as builder

ARG PROJECT_PATH="."
ARG DIST_PATH="./mokk"

# Create app directory
WORKDIR /workspace

# Copy source files
COPY . .

# Build app
RUN --mount=type=cache,target=/root/.cache/go-build go build -o ${DIST_PATH} ${PROJECT_PATH}

FROM alpine:latest

WORKDIR /app

COPY --from=builder /workspace/${DIST_PATH} ./

ENV SERVER_HOST 0.0.0.0

EXPOSE 80

CMD [ "./mokk", "start" ]
