conn = new Mongo();

const admin = conn.getDB('admin');
const { databases } = admin.adminCommand({ listDatabases: 1 });
const dbNames = databases.map(db => db.name);

if (!dbNames.includes(dbName)) {
  db = conn.getDB(dbName);

  db.createCollection(
    'users',
    {
      validator: {
        $jsonSchema: {
          bsonType: 'object',
          required: ['username', 'password'],
          properties: {
            username: {
              bsonType: 'string',
              minLength: 8,
              maxLength: 30,
              description: 'is required and must be a string with length between 8 and 30',
            },
            password: {
              bsonType: 'string',
              minLength: 8,
              maxLength: 30,
              description: 'is required and must be a string with length between 8 and 30',
            },
          },
        },
      },
    },
  );

  db.createCollection(
    'followers',
    {
      validator: {
        $jsonSchema: {
          bsonType: 'object',
          required: ['followerUserID', 'followeeUserID'],
          properties: {
            followerUserID: {
              bsonType: 'string',
              description: 'user ID of the follower; references the _id of a user in the "users" collection',
            },
            followeeUserID: {
              bsonType: 'string',
              description: 'user ID of the person being followed; references the _id of a user in the "users" collection',
            },
          },
        },
      },
    },
  );

  db.createCollection(
    'tweets',
    {
      validator: {
        $jsonSchema: {
          bsonType: 'object',
          required: ['userID', 'username', 'text'],
          properties: {
            userID: {
              bsonType: 'string',
              description: 'references the _id of a user in the "users" collection',
            },
            username: {
              bsonType: 'string',
              description: 'is the username of the user with the given userID',
            },
            text: {
              bsonType: 'string',
              minLength: 1,
              maxLength: 100,
              description: 'is required and must be a string with length between 1 and 100',
            },
          },
        },
      },
    },
  );
}
