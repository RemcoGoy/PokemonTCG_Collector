import 'package:application/main.dart';
import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'HttpService.dart';

void signUserIn(String email, String password, BuildContext context) async {
  if(email.isEmpty || password.isEmpty){
    print('password or email not given');
  }else{
    final userStorage = Hive.box('userStorage');

    final response = await post('/login', {
      email: email,
      password: password
    });

    if(response.statusCode == 200){
      print(response.body);
    }else{
      print("Something whent wrong: " + response.body);
    }
    // userStorage.put('email', email);
    // Navigator.push(context, MaterialPageRoute(builder: (context) => MainApp()));
  }
}

void signUserUp(String email, String password, String username, BuildContext context) {
  if(email.isEmpty || password.isEmpty){
    print('password or email not given');
  }else{

    Navigator.push(context, MaterialPageRoute(builder: (context) => MainApp()));
  }
}