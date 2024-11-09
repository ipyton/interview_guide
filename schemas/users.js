/* global use, db */
// MongoDB Playground
// Use Ctrl+Space inside a snippet or a string literal to trigger completions.

// The current database to use.
use('interview_guide');

// Create a new index in the collection.
db.getCollection('user_info')
    .createIndex(
        {
            open_id: 1
        }, {
        unique: true
    }
    );
db.getCollection('user_status')
    .createIndex(
        {
            open_id: 1
        }, {
        unique: true
    }
    );
