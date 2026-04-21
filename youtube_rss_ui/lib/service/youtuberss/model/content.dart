class Content {
  final int id;
  final String url;
  final int channel;
  final String title;
  final String image;
  final DateTime date;

  Content({
    required this.id,
    required this.url,
    required this.channel,
    required this.title,
    required this.image,
    required this.date,
  });

  factory Content.fromJson(Map<String, dynamic> json) {
    return Content(
      id: json['id'] ?? 0,
      url: json['url'] ?? '',
      channel: json['channel'] ?? 0,
      title: json['title'] ?? '',
      image: json['image'] ?? '',
      date: json['date'] != null 
          ? DateTime.parse(json['date']) 
          : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'url': url,
      'channel': channel,
      'title': title,
      'image': image,
      'date': date.toIso8601String(),
    };
  }
}