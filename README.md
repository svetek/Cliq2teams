# Tool for migrate from Zoho cliq to Micosoft Teams
## build
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../bin/cliq2teams-linux main
```

## USAGE 
Set variables in .env file: 
```env
tenantID="*****"
clientID="****"
clientSecret="****"

TeamName="Any team name"
TeamDescription="Some describtion"
TeamCreateDate="2015-03-14T11:22:17.043Z"
GuestAzObjectID="*****-f2d0-4dcd-a8b6-ff8a5b50cc4a"

parallelImportMessages=30
```

## create files/import directory and put export Zoho Cliq files
```bash
mkdir -p files/output && mkdir -p files/import
```

