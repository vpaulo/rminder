# Project rminder

Keep track of your tasks

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```

# Deploy

Build debian packages for rminder and rminder-caddy:
```
make package
```

Create `rminder` user on the host:
```
adduser --system --no-create-home --disabled-password --disabled-login rminder
```

Install packages on the host:
```
deb -i rminder.deb
deb -i rminder-caddy.deb
``` 

Enable systemd services:
```
systemctl daemon-reload

systemctl enable rminder
systemctl enable rminder-caddy
```

Start the services:
```
systemctl start rminder
systemctl status rminder

systemctl start rminder-caddy
systemctl status rminder-caddy
```

Check the logs:
```
journalctl -u rminder.service -f
journalctl -u rminder-caddy.service -f
```

# Testing Stripe payments

Ensure you have these variables in the `.env` file:
```
STRIPE_RETURN_TO_URL=http://localhost:3000
STRIPE_PREMIUM_PRICE_ID=...
STRIPE_SECRET_KEY=...
STRIPE_WEBHOOK_SECRET=...
```

`STRIPE_SECRET_KEY` you can get from your Stripe account. You will have to create a product and set a price and then copy the price id to `STRIPE_PREMIUM_PRICE_ID`.

Download Stripe CLI and login:

```
stripe login
```

Then run:

```
stripe listen --forward-to localhost:3000/post-checkout/webhook
```

Copy webhook secret and put it in `STRIPE_WEBHOOK_SECRET`.

You can now do test payments with Stripe!