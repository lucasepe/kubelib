#!/bin/bash

API_DIR=./apis

# Generate DeepCopy methods
controller-gen paths=${API_DIR}/... object
