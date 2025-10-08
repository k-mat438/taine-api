# Makefile for Database Operations

# 変数の定義
DB_CONTAINER = taine-api-db-1
DB_NAME = taine
DB_USER = postgres

# デフォルトターゲット
.DEFAULT_GOAL := help

# ヘルプ表示
.PHONY: help
help:
	@echo "利用可能なコマンド:"
	@echo "  migrate-up   - マイグレーションを実行"
	@echo "  migrate-down - マイグレーションを適用解除"
	@echo "  seed        - シードデータを投入"
	@echo "  reset       - データベースをリセット"

# マイグレーション
.PHONY: migrate-up
migrate-up:
	@echo "マイグレーションを実行中..."
	@for file in db/migrations/*_up.sql; do \
		echo "実行中: $$file"; \
		docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) < $$file; \
	done
	@echo "マイグレーションが完了しました"

# マイグレーション適用解除
.PHONY: migrate-down
migrate-down:
	@echo "マイグレーションを適用解除中..."
	@for file in db/migrations/*_down.sql; do \
		echo "実行中: $$file"; \
		docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) < $$file; \
	done
	@echo "マイグレーションの適用解除が完了しました"

# シードデータ
.PHONY: seed
seed:
	@echo "シードデータを投入中..."
	@for file in db/seeds/*.sql; do \
		echo "実行中: $$file"; \
		docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) < $$file; \
	done
	@echo "シードデータの投入が完了しました"

# データベースリセット
.PHONY: reset
reset:
	@echo "データベースをリセット中..."
	docker-compose down -v
	docker-compose up -d
	@echo "データベースの起動を待機中..."
	@sleep 3
	@$(MAKE) migrate-up
	@$(MAKE) seed
	@echo "データベースのリセットが完了しました"
