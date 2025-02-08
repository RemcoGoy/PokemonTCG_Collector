import 'package:application/components/Button.dart';
import 'package:application/components/Textfield.dart';
import 'package:application/main.dart';
import 'package:application/pages/Home.dart';
import 'package:application/pages/login.dart';
import 'package:application/services/AuthService.dart';
import 'package:flutter/material.dart';

class Signup extends StatelessWidget {
  Signup({super.key});

  final emailController = TextEditingController();
  final passwordController = TextEditingController();
  final usernameController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[300],
      body: SafeArea(
        child: Center(
          child: Column(
            children: [
              const SizedBox(height: 50),
              Icon(
                Icons.person,
                size: 100,
              ),
              const SizedBox(height: 50),
              Text(
                'Login',
                style: TextStyle(
                  color: Colors.grey[700],
                  fontSize: 18
                ),
              ),
              const SizedBox(height: 25),

              TCGTextField(
                controller: usernameController,
                hintText: 'Email',
              ),

              const SizedBox(height: 25),

              TCGTextField(
                controller: usernameController,
                hintText: 'Username',
              ),

              const SizedBox(height: 25),

              TCGTextField(
                controller: passwordController,
                hintText: 'Password',
                obscureText: true,
              ),

              const SizedBox(height: 25),

              TCGButton(
                buttonText: 'Sign up',
                onTap: () => signUserUp(emailController.text, passwordController.text, usernameController.text, context),
              ),
              
              const SizedBox(height: 25),

              GestureDetector(
                onTap: () => Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => MainApp())
                ),
                child: Text('Go back to login')
              )
            ],
          ),
        ),
      ),
    );
  }
}