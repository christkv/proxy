var MongoClient = require('mongodb').MongoClient
  , assert = require('assert');

MongoClient.connect('mongodb://localhost:50000/test?maxPoolSize=1', function(err, db) {
  assert.equal(null, err);
  assert.ok(db != null);

  db.collection('test1').insert({a:1}, function(err, r) {
    assert.equal(null, err);

    db.collection('test1').findOne({a:1}, function(err, doc) {
      assert.equal(null, err);
      assert.equal(1, doc.a);
      db.close();
    });
  });
});