import 'package:application/components/Button.dart';
import 'package:application/components/Textfield.dart';
import 'package:application/pages/SignUp.dart';
import 'package:application/services/AuthService.dart';
import 'package:flutter/material.dart';

class Login extends StatefulWidget {
  const Login({super.key});

  @override
  State<Login> createState() => _LoginState();
}

class _LoginState extends State<Login> {
  final emailController = TextEditingController();
  final passwordController = TextEditingController();
  
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
                controller: emailController,
                hintText: 'Email',
              ),

              const SizedBox(height: 25),

              TCGTextField(
                controller: passwordController,
                hintText: 'Password',
                obscureText: true,
              ),

              const SizedBox(height: 25),

              TCGButton(
                buttonText: 'Sign in',
                onTap: () => signUserIn(emailController.text, passwordController.text, context),
              ),

              const SizedBox(height: 25),

              GestureDetector(
                onTap: () => Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => Signup())
                ),
                child: Text('sign up')
              )
            ],
          ),
        ),
      ),
    );
  }
}