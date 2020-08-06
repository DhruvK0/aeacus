ablp() {
  echo "Building aeacus..."
  go build -ldflags '-s -w' -o ./aeacus .
  echo "Linux aeacus build successful!"

  echo "Building phocus..."
  go build -ldflags '-s -w' -tags phocus -o ./phocus .
  echo "Linux phocus build successful!"
}

abwp() {
  echo "Building aeacus..."
  GOOS=windows go build -ldflags '-s -w' -o ./aeacus.exe .
  echo "Windows aeacus build successful!"

  echo "Building phocus..."
  GOOS=windows go build -ldflags '-s -w' -tags phocus -o ./phocus.exe .
  echo "Windows phocus build successful!"
}

ablp
abwp
