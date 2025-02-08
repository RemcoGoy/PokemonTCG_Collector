import 'package:flutter/material.dart';

class TCGTextField extends StatelessWidget {
  final controller;
  final String hintText;
  final bool obscureText;

  const TCGTextField({
    super.key,
    required this.controller,
    required this.hintText,
    this.obscureText = false
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 25.0),
      child: TextField(
        controller: this.controller,
        obscureText: this.obscureText,
        decoration: InputDecoration(
          enabledBorder: OutlineInputBorder(
            borderSide: BorderSide(color: Colors.deepPurple)
          ),
          focusedBorder: OutlineInputBorder(
            borderSide: BorderSide(color: Colors.deepPurpleAccent)
          ),
          fillColor: Colors.grey.shade100,
          filled: true,
          hintText: this.hintText,
          hintStyle: TextStyle(
            color: Colors.grey[700]
          )
        ),
      ),
    );
  }
}