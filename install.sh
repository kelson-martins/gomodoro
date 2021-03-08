user="$SUDO_USER"

if [ "X$user" = "X" ]; then
  echo "the setup script must be run with sudo"
  exit 1
fi

mkdir -p ${HOME}/gomodoro
cp config/config.yaml ${HOME}/gomodoro
chown -R $user.$user ${HOME}/gomodoro

mkdir -p /etc/gomodoro/
cp assets/tone.mp3 /etc/gomodoro