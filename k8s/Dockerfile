############################
# STEP 2 build a small image
############################
FROM alpine:3.7
# Copy our static executable
COPY ./makako-api /makako-api/bin/makako-api

# Run the binary
ENTRYPOINT ["/makako-api/bin/makako-api", "-deployed=true"]
EXPOSE 7777
