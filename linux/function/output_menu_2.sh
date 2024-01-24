echo "Which Operating System do you like?"

select os in Ubuntu LinuxMint Windows8 Windows7 WindowsXP
do
  case $os in
    "Ubuntu"|"LinuxMint")
      echo "I also use $os."
      break
    ;;
    "Windows8" | "Windows7" | "WindowsXP")
      echo "Why don't you try Linux?"
      break
    ;;
    *)
      echo "Invalid entry."
      break
    ;;
  esac
done