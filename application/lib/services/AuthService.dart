import 'dart:convert';

import 'package:application/main.dart';
import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'HttpService.dart';

void signUserIn(String email, String password, BuildContext context) async {
  if(email.isEmpty || password.isEmpty){
    print('password or email not given');
  }else{
    final userStorage = Hive.box('userStorage');

    final response = await post('auth/login', {
      "email": email,
      "password": password
    });

    if(response.statusCode == 200){
      userStorage.put('email', email);
      userStorage.put('AuthToken', json.decode(response.body)['token']);
    }else{
      print("Something whent wrong: " + response.body);
    }
    Navigator.push(context, MaterialPageRoute(builder: (context) => MainApp()));
  }
}

void signUserUp(String email, String password, String username, BuildContext context) async {
  if(email.isEmpty || password.isEmpty){
    print('password or email not given');
  }else{
    var body = {
      "email": email,
      "password": password,
      "username": username
    };

    print(body);
    final response = await post('auth/signup', body);

    if(response.statusCode == 200){
      print(response.body);
    }else{
      print("Something whent wrong: " + response.body);
    }
    
    // Navigator.push(context, MaterialPageRoute(builder: (context) => MainApp()));
  }
}

void logout(BuildContext context){
  final userStorage = Hive.box('userStorage');
  userStorage.clear();

  Navigator.push(context, MaterialPageRoute(builder: (context) => MainApp()));
}