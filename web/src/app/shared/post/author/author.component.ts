import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { AbstractPostComponent } from '../abstract-post.component';

@Component({
  selector: 'app-post-author',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './author.component.html',
  styleUrls: ['./author.component.scss'],
})
export class PostAuthorComponent extends AbstractPostComponent {

  /**
   * Use to indicate whether the display should show or not
   */
  @Input()
  displayName = true;

}
