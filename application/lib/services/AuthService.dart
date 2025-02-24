import 'dart:convert';

import 'package:application/main.dart';
import 'package:application/pages/Conformation.dart';
import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'HttpService.dart';

void signUserIn(String email, String password, BuildContext context) async {
  if(email.isEmpty || password.isEmpty){
    print('password or email not given');
  }else{
    final userStorage = Hive.box('userStorage');

    final response = await post('auth/login', body: {
      "email": email,
      "password": password
    });

    if(response.statusCode == 200){
      userStorage.put('email', email);
      userStorage.put('AuthToken', json.decode(response.body)['token'].toString());
      Navigator.push(context, MaterialPageRoute(builder: (context) => MainApp()));
    }else{
      print("Something whent wrong: " + response.body);
    }
  }
}

void signUserUp(String email, String password, String username, BuildContext context) async {
  if(email.isEmpty || password.isEmpty){
    print('password or email not given');
  }else{
    final userStorage = Hive.box('userStorage');

    final response = await post('auth/signup', body: {
      "email": email,
      "password": password,
      "username": username
    });

    if(response.statusCode == 200){
      Navigator.push(context, MaterialPageRoute(builder: (context) => Conformation()));
    }else{
      print("Something whent wrong: " + response.body);
    }
    
  }
}

void logout(BuildContext context) async {
  final userStorage = Hive.box('userStorage');
  final response = await authPost('auth/logout');
    userStorage.clear();
  
  if(response.statusCode == 200){
    Navigator.push(context, MaterialPageRoute(builder: (context) => MainApp()));
  }else{
    print("Something whent wrong: " + response.body);
  }
}