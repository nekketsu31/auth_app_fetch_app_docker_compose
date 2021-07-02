npm install
npx sequelize-cli db:migrate --config "config/config.json" --env "local" up
npm run start PORT=3000 HOST=localhost