# Tutorial to reset the node

### 1. Remove all runtime directories using the following command in the directory where you have the orai.env file

```bash
sudo rm -rf .oraid/
```

### 2. Pull and recreate the latest version of the orai image:

```
docker-compose pull orai && docker-compose up -d --force-recreate
```