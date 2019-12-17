import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { AbstractPostComponent } from '../abstract-post.component';

@Component({
  selector: 'app-post-categories',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './categories.component.html',
  styleUrls: ['./categories.component.scss'],
})
export class PostCategoriesComponent extends AbstractPostComponent implements OnInit {

  /**
   * List of categories
   */
  categories: Array<Category>;

  ngOnInit(): void {
    this.categories = this.post.categories;
  }

}
