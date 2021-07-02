const express 			= require('express');
const router 			= express();

const AuthController 	= require('./../controllers/AuthController');

const passport	= require('passport');

require('./../middleware/passport')(passport)
/* GET home page. */
router.get('/', function(req, res, next) {
  res.json({status:"success", message:"Parcel Pending API", data:{"version_number":"v1.0.0"}})
});

router.post('/register', AuthController.create);
router.post('/login', AuthController.login);
router.get('/users', passport.authenticate('jwt', {session:false}), AuthController.get);

router.get('')

module.exports = router;