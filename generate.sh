#!/usr/bin/env bash

protoc  greet/greetpb/greet.proto  --go_out=plugins=grpc:.

protoc  calculator/calculator_pb/calculator.proto  --go_out=plugins=grpc:.

