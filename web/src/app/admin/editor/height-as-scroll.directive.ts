import { AfterViewInit, Directive, ElementRef, HostListener } from '@angular/core';

@Directive({ selector: '[appHeightAsScroll]' })
export class HeightAsScrollDirective implements AfterViewInit {

  /**
   * Store a previous scroll height for comparing when keypress but still in the same line
   */
  private previousScrollHeight: number;

  constructor(private host: ElementRef) { }

  ngAfterViewInit(): void {
    setTimeout(() => this.resize());
  }

  @HostListener('change')
  onChange(): void {
    this.resize();
  }

  @HostListener('document:keypress', ['$event'])
  onKeyPress(event: KeyboardEvent): void {
    this.resize();
  }

  private resize(): void {
    const elem: HTMLElement = this.host.nativeElement as HTMLElement;

    if (elem.scrollHeight !== this.previousScrollHeight) {
      elem.style.height = elem.scrollHeight.toString() + 'px';
    }

    this.previousScrollHeight = elem.scrollHeight;
  }

}
