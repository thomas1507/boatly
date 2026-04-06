# Boatly

A Go + SQLite web application served by Caddy.

## 🌐 Live Websites

- **Production**: https://boatly.nl (en www.boatly.nl)
- **Test**: https://t.boatly.nl

## 📁 Project Structure

```
boatly/
├── cmd/server/          # Application entry point
│   └── main.go
├── internal/
│   ├── database/        # SQLite setup
│   │   └── db.go
│   └── handlers/        # HTTP handlers
│       └── handlers.go
├── .github/workflows/   # GitHub Actions deployment
│   ├── deploy.yml       # Production (main branch)
│   └── deploy-test.yml  # Test (test branch)
├── data/                # SQLite database storage
├── static/              # Static assets (CSS, JS, images)
├── boatly.service       # Production systemd service
├── boatly-test.service  # Test systemd service
└── go.mod
```

## 🚀 Automatische Deployment

Deze repository gebruikt **GitHub Actions** voor automatische deployment:

| Branch | Doel | URL | Poort |
|--------|------|-----|-------|
| `main` | Productie | https://boatly.nl | 3050 |
| `test` | Test omgeving | https://t.boatly.nl | 3051 |

### Hoe werkt het?

1. **Push naar `test` branch** → Auto-deploy naar test server (t.boatly.nl)
2. **Test je wijzigingen** op https://t.boatly.nl
3. **Merge test → main** via Pull Request
4. **Push naar `main`** → Auto-deploy naar productie (boatly.nl)

### Deployment Workflow

```bash
# 1. Werk op test branch
git checkout test
# ... maak je wijzigingen ...
git add .
git commit -m "Beschrijving van wijzigingen"
git push origin test

# 2. GitHub Actions deployt automatisch naar t.boatly.nl
# 3. Test op https://t.boatly.nl

# 4. Als alles werkt, merge naar main
git checkout main
git merge test
git push origin main

# 5. GitHub Actions deployt automatisch naar boatly.nl
```

## 🛠️ Server Setup (Eenmalig)

### 1. SSH Key Genereren (op server - indien nog niet gedaan)

```bash
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/github_actions_deploy
# Geen wachtwoord invoeren (druk 2x Enter)
cat ~/.ssh/github_actions_deploy.pub >> ~/.ssh/authorized_keys
cat ~/.ssh/github_actions_deploy
# Kopieer deze private key naar GitHub secrets!
```

### 2. GitHub Secrets Toevoegen

1. Ga naar: https://github.com/thomas1507/boatly/settings/secrets/actions
2. Voeg deze secrets toe:
   - **HOST**: `142.132.201.125`
   - **USERNAME**: `thomas`
   - **SSH_PRIVATE_KEY**: (je private key van stap 1)

### 3. Systemd Services Installeren

```bash
# Productie service
sudo cp /home/thomas/boatly/boatly.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable boatly
sudo systemctl start boatly

# Test service  
sudo cp /home/thomas/boatly/boatly-test.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable boatly-test
sudo systemctl start boatly-test
```

## 💻 Development

```bash
# Run locally
cd /home/thomas/boatly
go run cmd/server/main.go

# Build for production
GOOS=linux GOARCH=amd64 go build -o boatly-server cmd/server/main.go
```

## 🔄 Handmatige Commands (indien nodig)

```bash
# Check status
sudo systemctl status boatly
sudo systemctl status boatly-test

# Herstart services
sudo systemctl restart boatly
sudo systemctl restart boatly-test

# Bekijk logs
sudo journalctl -u boatly -f
sudo journalctl -u boatly-test -f
```

## 📝 Belangrijke Bestanden

- `cmd/server/main.go` - De hoofdapplicatie
- `internal/` - Interne packages
- `static/` - Frontend assets
- `data/` - SQLite database (niet in git!)
- `.github/workflows/` - Deployment configuratie

## 🆘 Problemen?

1. **GitHub Actions faalt**: Check of `SSH_PRIVATE_KEY`, `HOST`, en `USERNAME` correct zijn ingesteld
2. **Website niet bereikbaar**: Check `sudo systemctl status boatly`
3. **Wijzigingen niet zichtbaar**: Browser cache legen (Ctrl+Shift+R)

---

**Happy sailing!** ⚓🚢
