class Collection {
  final String name;
  final String id;

  Collection({
    required this.name,
    required this.id
  });

  factory Collection.fromJson(Map<String, dynamic> json) {
    return Collection(name: json['name'], id: json['id']);
  }
}