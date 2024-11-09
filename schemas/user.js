use('interview_guide');

db.user_status.deleteMany({})
db.user_info.deleteMany({})

use('interview_guide');

// Create a new index in the collection.
db.getCollection('user_info')
    .createIndex(
        {
            openid: 1
        }, {
        unique: true
    }
    );
db.getCollection('user_status')
    .createIndex(
        {
            openid: 1
        }, {
        unique: true
    }
    );
