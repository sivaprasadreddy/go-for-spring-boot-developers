package com.sivalabs.bookmarks.domain;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.jdbc.JdbcTest;
import org.springframework.jdbc.core.simple.JdbcClient;
import org.springframework.test.context.jdbc.Sql;

import java.time.Instant;
import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;

@JdbcTest(properties = {
   "spring.test.database.replace=none",
   "spring.datasource.url=jdbc:tc:postgresql:16-alpine:///db"
})
@Sql("classpath:/test_data.sql")
class BookmarkRepositoryTest {

    @Autowired
    JdbcClient jdbcClient;

    BookmarkRepository bookmarkRepository;

    @BeforeEach
    void setUp() {
        bookmarkRepository = new BookmarkRepository(jdbcClient);
    }

    @Test
    void shouldFindAllBookmarks() {
        List<Bookmark> bookmarks = bookmarkRepository.findAll();
        assertThat(bookmarks).isNotEmpty();
        assertThat(bookmarks).hasSize(6);
    }

    @Test
    void shouldCreateBookmark() {
        var bookmark = new Bookmark(null, "My Title", "https://sivalabs.in", Instant.now());
        var savedBookmark = bookmarkRepository.save(bookmark);
        assertThat(savedBookmark.id()).isNotNull();
    }

    @Test
    void shouldGetBookmarkById() {
        var bookmark = new Bookmark(null, "My Title", "https://sivalabs.in", Instant.now());
        var savedBookmark = bookmarkRepository.save(bookmark);

        var result = bookmarkRepository.findById(savedBookmark.id()).orElseThrow();
        assertThat(result.id()).isEqualTo(savedBookmark.id());
        assertThat(result.title()).isEqualTo(bookmark.title());
        assertThat(result.url()).isEqualTo(bookmark.url());
    }

    @Test
    void shouldEmptyWhenBookmarkNotFound() {
        Optional<Bookmark> bookmarkOptional = bookmarkRepository.findById(9999L);
        assertThat(bookmarkOptional).isEmpty();
    }

    @Test
    void shouldUpdateBookmark() {
        var bookmark = new Bookmark(null, "My Title", "https://sivalabs.in", Instant.now());
        var savedBookmark = bookmarkRepository.save(bookmark);

        var changedBookmark = new Bookmark(savedBookmark.id(), "My Updated Title", "https://www.sivalabs.in", bookmark.createdAt());
        bookmarkRepository.update(changedBookmark);

        var updatedBookmark = bookmarkRepository.findById(savedBookmark.id()).orElseThrow();
        assertThat(updatedBookmark.id()).isEqualTo(changedBookmark.id());
        assertThat(updatedBookmark.title()).isEqualTo(changedBookmark.title());
        assertThat(updatedBookmark.url()).isEqualTo(changedBookmark.url());
    }

    @Test
    void shouldDeleteBookmark() {
        var bookmark = new Bookmark(null, "My Title", "https://sivalabs.in", Instant.now());
        var savedBookmark = bookmarkRepository.save(bookmark);

        bookmarkRepository.delete(savedBookmark.id());

        Optional<Bookmark> optionalBookmark = bookmarkRepository.findById(savedBookmark.id());
        assertThat(optionalBookmark).isEmpty();
    }
}