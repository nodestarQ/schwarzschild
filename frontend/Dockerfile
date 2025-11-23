# Build stage
FROM oven/bun:latest AS builder

WORKDIR /app

# Copy package files
COPY package.json bun.lock ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Build the application
RUN bun run build

# Production stage
FROM oven/bun:latest

WORKDIR /app

# Copy built app from builder
COPY --from=builder /app/dist ./dist

# Copy package.json for reference
COPY package.json ./

EXPOSE 3547

# Serve the dist folder using bun
CMD ["bun", "run", "--", "bunx", "serve", "-s", "dist", "-l", "3000"]
