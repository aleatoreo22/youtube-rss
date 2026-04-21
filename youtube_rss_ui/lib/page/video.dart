import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:intl/intl.dart';
import 'package:youtube_rss_ui/service/youtuberss/model/content.dart';
import 'package:youtube_rss_ui/service/youtuberss/youtuberss.dart';

class VideoPage extends StatefulWidget {
  const VideoPage({super.key});

  @override
  State<VideoPage> createState() => _VideoPageState();
}

class _VideoPageState extends State<VideoPage> with TickerProviderStateMixin {
  final YoutubeRSS _api = YoutubeRSS();
  final ScrollController _scrollController = ScrollController();
  final int _userId = 1;
  final int _pageSize = 13;

  List<Content> _contents = [];
  bool _isLoading = true;
  bool _isLoadingMore = false;
  int _currentPage = 1;
  int _totalPages = 1;
  AnimationController? _controller;

  @override
  void initState() {
    super.initState();
    _loadContent();
    _scrollController.addListener(_onScroll);
    _controller =
        AnimationController(duration: const Duration(seconds: 2), vsync: this)
          ..addListener(() {
            setState(() {});
          });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    _controller?.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >=
        _scrollController.position.maxScrollExtent - 200) {
      if (!_isLoadingMore && _currentPage < _totalPages) {
        _loadMoreContent();
      }
    }
  }

  Future<void> _loadContent() async {
    setState(() => _isLoading = true);
    try {
      final pagination = await _api.getUserContentPaginated(
        _userId,
        1,
        _pageSize,
      );
      setState(() {
        _contents = pagination.items;
        _currentPage = pagination.page;
        _totalPages = pagination.totalPages;
        _isLoading = false;
      });
    } catch (e) {
      setState(() => _isLoading = false);
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Erro ao carregar: $e')));
      }
    }
  }

  Future<void> _loadMoreContent() async {
    setState(() => _isLoadingMore = true);
    try {
      final nextPage = _currentPage + 1;
      if (nextPage > _totalPages) return;

      final pagination = await _api.getUserContentPaginated(
        _userId,
        nextPage,
        _pageSize,
      );
      setState(() {
        _contents.addAll(pagination.items);
        _currentPage = pagination.page;
        _totalPages = pagination.totalPages;
        _isLoadingMore = false;
      });
    } catch (e) {
      setState(() => _isLoadingMore = false);
    }
  }

  Future<void> _updateContent() async {
    setState(() => _isLoading = true);

    await Future.delayed(Duration(seconds: 1));

    _controller?.forward(from: 0.0);
    await _controller?.reverse();

    try {
      final pagination = await _api.getUserContentPaginated(
        _userId,
        1,
        _pageSize,
      );
      setState(() {
        _contents = pagination.items;
        _currentPage = pagination.page;
        _totalPages = pagination.totalPages;
        _isLoading = false;
      });
    } catch (e) {
      setState(() => _isLoading = false);
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Erro ao carregar: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _contents.isEmpty
          ? const Center(child: Text('Nenhum conteúdo encontrado'))
          : ListView.builder(
              controller: _scrollController,
              itemCount: _contents.length + (_isLoadingMore ? 1 : 0),
              itemBuilder: (context, index) {
                if (index == _contents.length) {
                  return const Center(
                    child: Padding(
                      padding: EdgeInsets.all(16),
                      child: CircularProgressIndicator(),
                    ),
                  );
                }
                final content = _contents[index];
                return _buildContentCard(content);
              },
            ),
      floatingActionButton: FloatingActionButton(
        onPressed: _updateContent,
        tooltip: 'Atualizar',
        child: AnimatedBuilder(
          animation: _controller!,
          builder: (context, child) {
            return Icon(Icons.refresh, color: Colors.black);
          },
        ),
      ),
    );
  }

  Widget _buildContentCard(Content content) {
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      child: InkWell(
        onTap: () async {
          await openUrl(content.url);
        },
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            ClipRRect(
              borderRadius: const BorderRadius.vertical(
                top: Radius.circular(8),
              ),
              child: Image.network(
                content.image.replaceFirst("hqdefault", "maxresdefault"),
                width:
                    double.infinity, // Torna a imagem toda a largura disponível
                fit: BoxFit.fill,
                errorBuilder: (context, error, stackTrace) =>
                    const Icon(Icons.video_library),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    content.title,
                    style: const TextStyle(fontWeight: FontWeight.w600),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 4),
                  Text(
                    formatDate(content.date),
                    style: const TextStyle(color: Colors.grey),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

String formatDate(DateTime date) {
  return DateFormat('dd/MM/yyyy HH:mm').format(date);
}

Future<void> openUrl(String url) async {
  final uri = Uri.parse(url.replaceFirst('www.youtube.com', 'm.youtube.com'));
  if (await canLaunchUrl(uri)) {
    await launchUrl(uri, mode: LaunchMode.externalApplication);
  }
}
