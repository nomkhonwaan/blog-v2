import { Component, EventEmitter, HostListener, Output, ElementRef, ViewChild } from '@angular/core';

@Component({
  selector: 'app-more',
  templateUrl: './more.component.html',
  styleUrls: ['./more.component.scss'],
})
export class MoreComponent {

  /**
   * For emitting on scrolling event
   */
  @Output()
  scrolled: EventEmitter<number> = new EventEmitter();

  /**
   * An anchor element for detecting when scroll to the end
   */
  @ViewChild('anchor', { static: true })
  private anchor: ElementRef<HTMLElement>;

  /**
   * An intersection observer implementation
   */
  private intersectionObserver: IntersectionObserver;

  constructor(private host: ElementRef) { }

  @HostListener('window:scroll', ['$event'])
  onWindowScroll(event: Event): void {
    const scrollTop: number = event.target['scrollingElement'].scrollTop;

    this.scrolled.emit(scrollTop);
  }


}
