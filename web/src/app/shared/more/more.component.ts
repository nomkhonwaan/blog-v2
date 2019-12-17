import { Component, ElementRef, EventEmitter, OnInit, Output, ViewChild, OnDestroy } from '@angular/core';

@Component({
  selector: 'app-more',
  templateUrl: './more.component.html',
  styleUrls: ['./more.component.scss'],
})
export class MoreComponent implements OnInit, OnDestroy {

  /**
   * For emitting on scrolling event
   */
  @Output()
  scrolled: EventEmitter<null> = new EventEmitter(null);

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

  ngOnInit(): void {
    this.intersectionObserver = new IntersectionObserver(([entry]: Array<IntersectionObserverEntry>): void => {
      entry.isIntersecting && this.scrolled.emit();
    }, { root: null });

    this.intersectionObserver.observe(this.anchor.nativeElement);
  }

  ngOnDestroy(): void {
    this.intersectionObserver.disconnect();
  }

}
