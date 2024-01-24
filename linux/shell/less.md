# less

less -N 

/pattern  *  Search forward for (N-th) matching line.
?pattern  *  Search backward for (N-th) matching line.
&pattern  *  Display only matching lines.
g         *  Go to first line in file (or line N).
G         *  Go to last line in file (or line N).

f - 向前移动一屏
b - 向后移动一屏
d - 向前移动半屏
u - 向后移动半屏
F - tail -f, ctrl + c退出