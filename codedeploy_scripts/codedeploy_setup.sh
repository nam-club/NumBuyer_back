#! /bin/sh
sudo yum install -y docker
sudo systemctl start docker
sudo usermod -a -G docker ec2-user
sudo curl -L https://github.com/docker/compose/releases/download/1.28.5/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose -f /home/ec2-user/docker-compose.prd.yml down
docker-compose -f /home/ec2-user/docker-compose.prd.yml up -d