package repository

const (
	QueryCreateUser = `
		INSERT INTO users (id, username, password_hash) 
		VALUES (gen_random_uuid(), $1, $2)`

	QueryGetUserByLogin = `
		SELECT id, password_hash FROM users WHERE username = $1`

	getUserCoinsQuery = `
		SELECT coins
		FROM users
		WHERE id = $1;`

	getCoinHistoryQuery = `
    SELECT 
        u_from.username AS from_user, 
        u_to.username AS to_user, 
        t.amount
    FROM transactions t
    JOIN users u_from ON t.from_user = u_from.id
    JOIN users u_to ON t.to_user = u_to.id
    WHERE t.from_user = $1 OR t.to_user = $1
    ORDER BY t.id DESC;`

	QueryGetUserInventory = `
	SELECT type, quantity 
	FROM inventories 
	WHERE user_id = $1`

	QuerySubtractUserCoins = `
	UPDATE users 
	SET coins = coins - $1 
	WHERE id = $2 AND coins >= $1
	RETURNING coins;`

	QueryAddOrUpdateItemInInventory = `
    INSERT INTO inventories (user_id, type, quantity)
    VALUES ($1, $2, $3)
    ON CONFLICT (user_id, type) 
    DO UPDATE SET quantity = inventories.quantity + $3;
    `
	QueryInsertTransaction = `
	INSERT INTO transactions (from_user, to_user, amount)
	VALUES ($1, $2, $3)`

	QueryAddUserCoins = `
	UPDATE users 
	SET coins = coins + $1 
	WHERE id = $2
	RETURNING coins;`
)
