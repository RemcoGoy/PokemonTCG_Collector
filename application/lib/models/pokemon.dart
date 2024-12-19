class Pokemon {
  final String name;
  final String id;
  final String imageUrl;
  final String setId;

  Pokemon({
    required this.name,
    required this.id,
    required this.imageUrl,
    required this.setId
  });

  factory Pokemon.fromJson(Map<String, dynamic> json) {
    return Pokemon(name: json['name'], id: json['id'], imageUrl: json['images']['small'], setId: json['set']['id']);
  }
}