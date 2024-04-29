db = db.getSiblingDB('dungeondb');

db.createUser({
  user: 'gameuser',
  pwd: '$GAME_USER_PASS',
  roles: [{ role: 'readWrite', db: 'dungeondb' }],
});

// Create collections if they don't exist
const collections = ['players', 'users', 'lobbies', 'matches'];
collections.forEach(col => {
  if (!db.getCollectionNames().includes(col)) {
    db.createCollection(col);
  }
});
