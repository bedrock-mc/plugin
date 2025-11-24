# go
go mod tidy

# nodejs
cd packages/node
npm install --no-audit --no-fund
npm run build
cd ../..
