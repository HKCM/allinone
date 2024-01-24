#!/usr/bin/env bash
read -n1 -p "Do you want to continue [Y/N]? " answer 
case $answer in
  Y | y)  
    echo ""
    echo "fine, continue on...";; 
  N | n)  
    echo ""
    echo OK, goodbye
    exit;;
  *)  
    echo ""
    echo "Error Input..."
    exit;;
  esac
echo "This is the end of the script"