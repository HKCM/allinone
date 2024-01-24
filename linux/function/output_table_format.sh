#!/usr/bin/env bash

# Function to print a horizontal line
print_horizontal_line() {
  local columns=$1
  local line=""
  for ((i=1; i<=columns; i++)); do
    line+="------------+"
  done
  printf "+%s\n" "$line"
}

# Function to print table header
print_header() {
  local header_length=${#header_fields[@]}
  local header=""
  for ((i=1; i<=header_length; i++)); do
    header+=" %-10s |"
  done
  printf "|$header\n" "$@"
  print_horizontal_line $header_length
}

# Function to print table row
print_row() {
  local data_length=$#
  local row=""
  for ((i=1; i<=data_length; i++)); do
    row+=" %-10s |"
  done
  printf "|$row\n" "$@"
}


# Table header fields
header_fields=("Header 1" "Header 2" "Header 3" "Header 4")
data_field_1=("data 1" "data 2" "data 3" "da")
data_field_2=("data 4" "data 5" "data 6" "1")

# Print table
print_horizontal_line "${#header_fields[@]}"
print_header "${header_fields[@]}"
print_row "${data_field_1[@]}"
print_row "${data_field_2[@]}"
print_horizontal_line "${#header_fields[@]}"
