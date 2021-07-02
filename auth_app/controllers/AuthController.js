const User 			= require('../models').User;
const validator     = require('validator');

const getUniqueKeyFromBody = function(body){// this is so they can send in 3 options unique_key, email, or phone and it will work
    let unique_key = body.unique_key;
    if(typeof unique_key==='undefined'){
        if(typeof body.name != 'undefined'){
            unique_key = body.name
        }else if(typeof body.phone != 'undefined'){
            unique_key = body.phone
        }else{
            unique_key = null;
        }
    }

    return unique_key;
}
module.exports.getUniqueKeyFromBody = getUniqueKeyFromBody;

const createUser = async function(userInfo){
    let unique_key, auth_info, err;

    auth_info={}
    auth_info.status='create';
    
    unique_key = getUniqueKeyFromBody(userInfo);
    if(!unique_key) TE('An username or phone number was not entered.');
    
    if(validator.isAlphanumeric(unique_key)){//checks if only phone number was sent
        auth_info.method = 'phone';
        // userInfo.phone = unique_key;
        
        [err, user] = await to(User.create(userInfo));
        if(err) TE(err.message);

        return user;
    }else{
        TE('A valid username or phone number was not entered.');
    }
}
module.exports.createUser = createUser;


const authUser = async function(userInfo){//returns token
    let unique_key;
    let auth_info = {};
    auth_info.status = 'login';
    unique_key = getUniqueKeyFromBody(userInfo);
    if(!unique_key) TE('Please enter an phone number to login');
    if(!userInfo.password) TE('Please enter a password to login');
    
    let user;
    if(validator.isAlphanumeric(unique_key)){
        auth_info.method='phone';
        [err, user] = await to(User.findOne({where:{phone:userInfo.phone}}));
        if(err) TE(err.message);
    }else{
        TE('A valid phone number was not entered');
    }

    if(!user) TE('Not registered');
    [err, user] = await to(user.comparePassword(userInfo.password));
    if(err) TE(err.message);

    return user;

}
module.exports.authUser = authUser;


const create = async function(req, res){
    res.setHeader('Content-Type', 'application/json');
    const body = req.body;
    var generatePassword = await to(password_generate())
    body.password = generatePassword[1]
    
    if(!body.unique_key && !body.phone ){
        return ReE(res, 'Please enter an phone number to register.');
    } else if(!body.password){
        return ReE(res, 'Please enter a password to register.');
    }else{
        let err, user;
        
        [err, user] = await to(createUser(body));
        if(err) return ReE(res, err, 422);
        user.password = body.password
        return ReS(res, {message:'Successfully created new user.', user:user.toWeb(), token:user.getJWT()}, 201);
    }
}
module.exports.create = create;

const login = async function(req, res){
    let err, user;

    [err, user] = await to(authUser(req.body));
    if(err) return ReE(res, err, 422);
    let users = await user.toWeb()
    delete users.password

    return ReS(res, {token:user.getJWT(), user:users});
}
module.exports.login = login;

const get = async function(req, res){
    res.setHeader('Content-Type', 'application/json');
    let user = req.user;

    let users = await user.toWeb()
    delete users.password

    return ReS(res, {user:users});
}
module.exports.get = get;

const password_generate = async function(params) {
    var result           = '';
   var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
   var charactersLength = characters.length;
   for ( var i = 0; i < 4; i++ ) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
   }
   return result;
}
module.exports.password_generate = password_generate;