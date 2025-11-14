import { Component, Input, Output, EventEmitter, OnChanges, SimpleChanges } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { NoteService, Note, Content } from '../../services/note-service';

@Component({
  selector: 'app-note-editor-side-panel',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './note-editor-side-panel.html',
  styleUrl: './note-editor-side-panel.scss',
})
export class NoteEditorSidePanelComponent implements OnChanges {
  @Input() noteId: string | null = null;
  @Output() closePanel = new EventEmitter<void>();

  note: Note | null = null;

  constructor(private noteService: NoteService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['noteId'] && this.noteId) {
      this.noteService.getNoteById(this.noteId).subscribe({
        next: (note) => {
          this.note = note;
          console.log('Loaded note', note);
        },
        error: (err) => {
          console.error('Error fetching note', err);
          this.note = null;
        },
      });
    }
  }

  onClose(): void {
    this.closePanel.emit();
  }

  onContentInput(event: Event): void {
    // No debouncing needed, updates are triggered by blur or enter key
  }

  onContentBlur(): void {
    const paragraph = window.getSelection()?.anchorNode?.parentElement;
    if (paragraph && paragraph.tagName === 'P') {
      const contentId = paragraph.getAttribute('data-content-id');
      setTimeout(() => {
        if (contentId) {
          this.updateContent(contentId, paragraph.innerText);
        }
      }, 0);
    }
  }

  onContentKeydown(event: KeyboardEvent): void {
    if (event.key === 'Enter') {
      event.preventDefault();

      const selection = window.getSelection();
      if (!selection || !selection.rangeCount) return;

      const range = selection.getRangeAt(0);
      const { startContainer, startOffset } = range;

      let currentParagraph: Node | null = startContainer;
      if (currentParagraph.nodeType === Node.TEXT_NODE) {
        currentParagraph = currentParagraph.parentElement;
      }

      if (!currentParagraph || (currentParagraph as HTMLElement).tagName !== 'P' || !this.note) return;
      
      const contentId = (currentParagraph as HTMLElement).getAttribute('data-content-id');
      if (!contentId) return;

      const currentIndex = this.note.contents.findIndex(c => c.id === contentId);
      if (currentIndex === -1) return;

      const originalText = (currentParagraph as HTMLElement).textContent || '';
      const textBeforeCursor = originalText.substring(0, startOffset);
      const textAfterCursor = originalText.substring(startOffset);

      // Update current content
      console.log('Updating content ID:', contentId, 'with text:', textBeforeCursor);
      this.updateContent(contentId, textBeforeCursor);

      // Create new content
      const newContent: Content = {
        id: '', // Will be set by the backend
        noteID: this.note.id,
        data: textAfterCursor,
        type: 'text',
        version: 0,
        position: currentIndex + 1,
      };

      this.noteService.addContent(this.note.id, newContent, this.note.version).subscribe({
        next: (addedContent) => {
          newContent.id = addedContent.id;
          this.note?.contents.splice(currentIndex + 1, 0, newContent);
          if (this.note) {
            this.note.version++;
          }
          
          // Set focus on the new element after Angular renders it
          setTimeout(() => {
            const newParagraph = document.querySelector(`[data-content-id="${addedContent.id}"]`) as HTMLElement;
            if (newParagraph) {
              newParagraph.focus();
              const newRange = document.createRange();
              const newSelection = window.getSelection();
              newRange.setStart(newParagraph.childNodes[0] || newParagraph, 0);
              newRange.collapse(true);
              newSelection?.removeAllRanges();
              newSelection?.addRange(newRange);
            }
          }, 0);
        },
        error: (err) => {
          console.error('Error adding content', err);
          // Revert the change on error
          if (this.note) {
            this.note.contents[currentIndex].data = originalText;
          }
        },
      });
    } else if (event.key === 'Backspace') {
      const selection = window.getSelection();
      if (!selection || !selection.rangeCount) return;

      const range = selection.getRangeAt(0);
      const { startContainer, startOffset } = range;

      if (startOffset === 0) {
        let currentParagraph: Node | null = startContainer;
        if (currentParagraph.nodeType === Node.TEXT_NODE) {
          currentParagraph = currentParagraph.parentElement;
        }

        if (!currentParagraph || (currentParagraph as HTMLElement).tagName !== 'P' || !this.note) return;

        const contentId = (currentParagraph as HTMLElement).getAttribute('data-content-id');
        if (!contentId) return;

        const currentIndex = this.note.contents.findIndex(c => c.id === contentId);
        if (currentIndex > 0) {
          event.preventDefault();
          const previousContent = this.note.contents[currentIndex - 1];
          const currentContent = this.note.contents[currentIndex];
          const mergedData = previousContent.data + currentContent.data;

          this.updateContent(previousContent.id, mergedData);
          this.noteService.deleteContent(this.note.id, currentContent.id, this.note.version, currentContent.version).subscribe({
            next: () => {
              if (this.note) {
                this.note.version++;
                this.note.contents.splice(currentIndex, 1);
                setTimeout(() => {
                  const prevParagraph = document.querySelector(`[data-content-id="${previousContent.id}"]`) as HTMLElement;
                  if (prevParagraph) {
                    prevParagraph.focus();
                    const newRange = document.createRange();
                    const newSelection = window.getSelection();
                    const textNode = prevParagraph.childNodes[0] || prevParagraph;
                    newRange.setStart(textNode, previousContent.data.length);
                    newRange.collapse(true);
                    newSelection?.removeAllRanges();
                    newSelection?.addRange(newRange);
                  }
                }, 0);
              }
            },
            error: (err) => {
              console.error('Error deleting content', err);
            }
          });
        }
      }
    }
  }

  private updateContent(contentId: string, newText: string): void {
    console.log('updateContent called for contentId:', contentId, 'newText:', newText);
    if (!this.note) {
      console.log('updateContent: Note is null, returning.');
      return;
    }

    const contentIndex = this.note.contents.findIndex(c => c.id === contentId);
    if (contentIndex === -1) {
      console.log('updateContent: Content not found, returning.');
      return;
    }

    const originalContent = this.note.contents[contentIndex];
    if (originalContent.data === newText) {
      console.log('updateContent: Content data is the same, returning.');
      return;
    }

    const updatedContent: Content = { ...originalContent, data: newText };
    console.log('Calling noteService.updateContent with:', updatedContent);

    this.noteService.updateContent(updatedContent).subscribe({
      next: () => {
        console.log('noteService.updateContent: next callback triggered.');
        if (this.note) {
          this.note.contents[contentIndex] = { ...updatedContent, version: originalContent.version + 1 };
          console.log(`Content ${contentId} updated successfully`);
        }
      },
      error: (err) => {
        console.error(`noteService.updateContent: Error updating content ${contentId}`, err);
        // Revert the UI change on error
        const paragraph = document.querySelector(`[data-content-id="${contentId}"]`);
        if (paragraph) {
          paragraph.textContent = originalContent.data;
        }
      },
    });
  }
}
