#!/bin/bash
# Run PHPStan linting via Docker

docker run --rm \
  -v "$(pwd):/app" \
  -w /app \
  --platform linux/amd64 \
  php:8.1-cli \
  bash -c "
    # Install composer if not present
    if [ ! -f /usr/local/bin/composer ]; then
      curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer
    fi
    
    # Install dependencies
    composer install --no-interaction --quiet
    
    # Run PHPStan
    php vendor/bin/phpstan analyze --no-progress
  "

