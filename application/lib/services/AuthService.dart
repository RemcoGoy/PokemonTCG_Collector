import 'package:application/main.dart';
import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';

void signUserIn(String email, String password, BuildContext context) {
  if(email.isEmpty || password.isEmpty){
    print('password or email not given');
  }else{
    final userStorage = Hive.box('userStorage');

    userStorage.put('email', email);
    Navigator.push(context, MaterialPageRoute(builder: (context) => MainApp()));
  }
}

void signUserUp(String email, String password, String username, BuildContext context) {
  if(email.isEmpty || password.isEmpty){
    print('password or email not given');
  }else{
    final userStorage = Hive.box('userStorage');

    userStorage.put('email', email);
    Navigator.push(context, MaterialPageRoute(builder: (context) => MainApp()));
  }
}