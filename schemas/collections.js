use('interview_guide');

db.collections.deleteMany({})
db.collection_items.deleteMany({})

use('interview_guide');

// Create a new index in the collection.
db.getCollection('collections')
    .createIndex(
        {
            openid: 1,
            collection_id:1
        }, {
            unique: true
        }
    );
// db.getCollection('collection_item')
//     .createIndex(
//         {
//             openid: 1
//         }, {
//             unique: true
//         }
//     );
