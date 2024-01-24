#!/usr/bin/env bash
SYSTEM=`uname  -s`
if [ $SYSTEM = "Linux" ] ; then
   echo "Linux"
elif
    [ $SYSTEM = "FreeBSD" ] ; then
   echo "FreeBSD"
elif
    [ $SYSTEM = "Solaris" ] ; then
    echo "Solaris"
elif
    [ $SYSTEM = "Darwin" ] ; then
    echo "Darwin"
else
    echo  "What?"
fi