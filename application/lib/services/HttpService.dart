import 'dart:convert';

import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:http/http.dart';
import 'package:http/http.dart' as http;

Future<Response> get(String uri) async {
  final serverUrl = dotenv.env['SERVER_URL'];
  if(serverUrl == null){
    throw Exception('serverUrl not found');
  }
  final url = Uri.parse(serverUrl + uri);
  print(url);
  return await http.get(url);
}

Future<Response> post(String uri, Object? body) async {
  final serverUrl = dotenv.env['SERVER_URL'];
  if(serverUrl == null){
    throw Exception('serverUrl not found');
  }
  final url = Uri.parse(serverUrl + uri);
  print(url);
  return await http.post(
    url,
    body: jsonEncode(body)
  );
}