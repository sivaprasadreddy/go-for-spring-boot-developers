package com.sivalabs.bookmarks.api;

import com.sivalabs.bookmarks.domain.Bookmark;
import com.sivalabs.bookmarks.domain.BookmarkRepository;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

import java.time.Instant;
import java.util.List;

@RestController
@RequestMapping("/api/bookmarks")
class BookmarksController {
    private final BookmarkRepository bookmarkRepository;

    BookmarksController(BookmarkRepository bookmarkRepository) {
        this.bookmarkRepository = bookmarkRepository;
    }

    @GetMapping
    List<Bookmark> findAll() {
        return bookmarkRepository.findAll();
    }

    @GetMapping("/{id}")
    ResponseEntity<Bookmark> findById(@PathVariable Long id) {
        return bookmarkRepository.findById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    Bookmark save(@RequestBody @Valid Bookmark bookmark) {
        var b = new Bookmark(null, bookmark.title(), bookmark.url(), Instant.now());
        return bookmarkRepository.save(b);
    }

    @PutMapping("/{id}")
    void update(@PathVariable Long id, @RequestBody @Valid Bookmark bookmark) {
        var b = new Bookmark(id, bookmark.title(), bookmark.url(), bookmark.createdAt());
        bookmarkRepository.update(b);
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    void delete(@PathVariable Long id) {
        bookmarkRepository.delete(id);
    }

}
